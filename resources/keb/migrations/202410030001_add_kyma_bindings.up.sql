CREATE TABLE IF NOT EXISTS bindings (
    id VARCHAR(255) NOT NULL,
    runtime_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    -- represents algorithm used to generate a kubeconfig, initialy: adminkubeconfig or tokenrequest
    type TEXT, 
    -- content of the kubeconfig
	kubeconfig TEXT, 
    -- expiration seconds
    expiration_seconds INTEGER,
    -- allow for the same binding id to be used for different runtimes
    PRIMARY KEY(id, runtime_id), 
    version integer NOT NULL
);
