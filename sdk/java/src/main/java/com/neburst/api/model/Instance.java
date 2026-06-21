package com.neburst.api.model;

import com.google.gson.annotations.SerializedName;

/**
 * Represents a compute instance.
 */
public class Instance {

    private String uuid;
    private String name;
    private String type;
    private String status;
    private String region;
    private String hostname;

    @SerializedName("pay_cycle")
    private String payCycle;

    @SerializedName("auto_renew")
    private boolean autoRenew;

    @SerializedName("next_pay_at")
    private String nextPayAt;

    @SerializedName("created_at")
    private String createdAt;

    public String getUuid() {
        return uuid;
    }

    public String getName() {
        return name;
    }

    public String getType() {
        return type;
    }

    public String getStatus() {
        return status;
    }

    public String getRegion() {
        return region;
    }

    public String getHostname() {
        return hostname;
    }

    public String getPayCycle() {
        return payCycle;
    }

    public boolean isAutoRenew() {
        return autoRenew;
    }

    public String getNextPayAt() {
        return nextPayAt;
    }

    public String getCreatedAt() {
        return createdAt;
    }

    @Override
    public String toString() {
        return "Instance{" +
                "uuid='" + uuid + '\'' +
                ", name='" + name + '\'' +
                ", type='" + type + '\'' +
                ", status='" + status + '\'' +
                ", region='" + region + '\'' +
                ", hostname='" + hostname + '\'' +
                ", payCycle='" + payCycle + '\'' +
                ", autoRenew=" + autoRenew +
                ", nextPayAt='" + nextPayAt + '\'' +
                ", createdAt='" + createdAt + '\'' +
                '}';
    }
}
