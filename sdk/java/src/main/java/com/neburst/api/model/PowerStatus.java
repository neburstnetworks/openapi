package com.neburst.api.model;

import com.google.gson.annotations.SerializedName;

/**
 * Represents the power status of a compute instance.
 */
public class PowerStatus {

    private String status;

    @SerializedName("is_installing")
    private boolean isInstalling;

    public String getStatus() {
        return status;
    }

    public boolean isInstalling() {
        return isInstalling;
    }

    @Override
    public String toString() {
        return "PowerStatus{" +
                "status='" + status + '\'' +
                ", isInstalling=" + isInstalling +
                '}';
    }
}
