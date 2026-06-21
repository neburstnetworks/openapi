package com.neburst.api;

import com.google.gson.reflect.TypeToken;
import com.neburst.api.model.Instance;
import com.neburst.api.model.PowerStatus;
import com.neburst.api.model.Traffic;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * API client for compute instance operations.
 */
public class ComputeApi {

    private final NeburstClient client;

    public ComputeApi(NeburstClient client) {
        this.client = client;
    }

    /**
     * Lists all compute instances.
     *
     * @return list of instances
     */
    public List<Instance> listInstances() throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/compute/instances", null, null,
                new TypeToken<List<Instance>>() {}.getType());
    }

    /**
     * Gets a single compute instance by ID.
     *
     * @param id the instance UUID
     * @return the instance details
     */
    public Instance getInstance(String id) throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/compute/instances/" + id, null, null,
                Instance.class);
    }

    /**
     * Gets the power status of an instance.
     *
     * @param id the instance UUID
     * @return the power status
     */
    public PowerStatus getInstanceStatus(String id) throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/compute/instances/" + id + "/status", null, null,
                PowerStatus.class);
    }

    /**
     * Gets the traffic information for an instance.
     *
     * @param id the instance UUID
     * @return the traffic info with package details
     */
    public Traffic getInstanceTraffic(String id) throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/compute/instances/" + id + "/traffic", null, null,
                Traffic.class);
    }

    /**
     * Performs a power action on an instance.
     *
     * @param id     the instance UUID
     * @param action the power action: "power-on", "power-off", "power-cycle", "power-reset"
     */
    public void powerAction(String id, String action) throws NeburstApiException, IOException, InterruptedException {
        Map<String, Object> body = new HashMap<>();
        body.put("action", action);
        client.doRequest("POST", "/open/v1/compute/instances/" + id + "/power", null, body);
    }

    /**
     * Rebuilds an instance with a new OS profile.
     *
     * @param id         the instance UUID
     * @param profileId  the OS profile ID
     * @param hostname   the new hostname (nullable)
     * @param publicKeys list of SSH public keys to inject (nullable)
     */
    public void rebuildInstance(String id, int profileId, String hostname, List<String> publicKeys)
            throws NeburstApiException, IOException, InterruptedException {
        Map<String, Object> body = new HashMap<>();
        body.put("profile_id", profileId);
        if (hostname != null) {
            body.put("hostname", hostname);
        }
        if (publicKeys != null) {
            body.put("public_keys", publicKeys);
        }
        client.doRequest("POST", "/open/v1/compute/instances/" + id + "/rebuild", null, body);
    }

    /**
     * Boots an instance into rescue mode with the specified profile.
     *
     * @param id        the instance UUID
     * @param profileId the rescue OS profile ID
     */
    public void rescueInstance(String id, int profileId) throws NeburstApiException, IOException, InterruptedException {
        Map<String, Object> body = new HashMap<>();
        body.put("profile_id", profileId);
        client.doRequest("POST", "/open/v1/compute/instances/" + id + "/rescue", null, body);
    }
}
