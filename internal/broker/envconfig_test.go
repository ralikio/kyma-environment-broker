package broker

import (
    "os"
    "testing"

	"github.com/vrischmann/envconfig"
    "github.com/stretchr/testify/assert"
)
type MyConfig struct {
    SapConveregedCloudMapping []struct {
        BtpRegion string
        SkrRegion string
    } `envconfig:""`
}

func TestEnvConfig(t *testing.T) {

    t.Run("should parse one to one mapping", func(t *testing.T) {
        // given 
        os.Setenv("APP_SAP_CONVEREGED_CLOUD_MAPPING", "{cf-eu20-staging,eu-de-1}")
        var cfg MyConfig

        // when 
        err := envconfig.InitWithPrefix(&cfg, "APP")
        
        // then
        assert.NoError(t, err)
        assert.NotNil(t, cfg.SapConveregedCloudMapping)
        assert.Equal(t, 1, len(cfg.SapConveregedCloudMapping))
        assert.Equal(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[0].BtpRegion)
        assert.Equal(t, "eu-de-1", cfg.SapConveregedCloudMapping[0].SkrRegion)

    })

    t.Run("should parse one to many mapping", func(t *testing.T) {
        // given 
        os.Setenv("APP_SAP_CONVEREGED_CLOUD_MAPPING", "{cf-eu20-staging,eu-de-1},{cf-eu20-staging,eu-de-2}")
        var cfg MyConfig

        // when 
        err := envconfig.InitWithPrefix(&cfg, "APP")
        
        // then
        assert.NoError(t, err)
        assert.NotNil(t, cfg.SapConveregedCloudMapping)
        assert.Equal(t, 2, len(cfg.SapConveregedCloudMapping))
        assert.Equal(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[0].BtpRegion)
        assert.Equal(t, "eu-de-1", cfg.SapConveregedCloudMapping[0].SkrRegion)

        assert.Equal(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[1].BtpRegion)
        assert.Equal(t, "eu-de-2", cfg.SapConveregedCloudMapping[1].SkrRegion)
    })

    t.Run("does not parse mappings with newlines", func(t *testing.T) {
        // given 
        os.Setenv("APP_SAP_CONVEREGED_CLOUD_MAPPING", `
            {cf-eu20-staging,eu-de-1},
            {cf-eu20-staging,eu-de-2}
        `)
        var cfg MyConfig

        // when 
        err := envconfig.InitWithPrefix(&cfg, "APP")
        
        // then
        assert.NoError(t, err)
        assert.NotNil(t, cfg.SapConveregedCloudMapping)
        assert.Equal(t, 2, len(cfg.SapConveregedCloudMapping))
        assert.NotEqual(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[0].BtpRegion)
        assert.Equal(t, "eu-de-1", cfg.SapConveregedCloudMapping[0].SkrRegion)

        assert.NotEqual(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[1].BtpRegion)
        assert.NotEqual(t, "eu-de-2", cfg.SapConveregedCloudMapping[1].SkrRegion)
    })


    t.Run("should parse multiple one to many mappings", func(t *testing.T) {
        // given 
        os.Setenv("APP_SAP_CONVEREGED_CLOUD_MAPPING", `{cf-eu20-staging,eu-de-1},{cf-eu20-staging,eu-de-2},{cf-eu21-staging,eu-pl-1},{cf-eu21-staging,eu-pl-2}`)
        var cfg MyConfig

        // when 
        err := envconfig.InitWithPrefix(&cfg, "APP")
        
        // then
        assert.NoError(t, err)
        assert.NotNil(t, cfg.SapConveregedCloudMapping)
        assert.Equal(t, 4, len(cfg.SapConveregedCloudMapping))

        assert.Equal(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[0].BtpRegion)
        assert.Equal(t, "eu-de-1", cfg.SapConveregedCloudMapping[0].SkrRegion)

        assert.Equal(t, "cf-eu20-staging", cfg.SapConveregedCloudMapping[1].BtpRegion)
        assert.Equal(t, "eu-de-2", cfg.SapConveregedCloudMapping[1].SkrRegion)
    
        assert.Equal(t, "cf-eu21-staging", cfg.SapConveregedCloudMapping[2].BtpRegion)
        assert.Equal(t, "eu-pl-1", cfg.SapConveregedCloudMapping[2].SkrRegion)

        assert.Equal(t, "cf-eu21-staging", cfg.SapConveregedCloudMapping[3].BtpRegion)
        assert.Equal(t, "eu-pl-2", cfg.SapConveregedCloudMapping[3].SkrRegion)
    })

}