package com.neburst.api.model;

import com.google.gson.annotations.SerializedName;

/**
 * Represents a single traffic package within an instance's traffic info.
 */
public class TrafficPackage {

    private String name;

    @SerializedName("capacity_gb")
    private int capacityGb;

    @SerializedName("used_gb")
    private double usedGb;

    @SerializedName("reset_cycle")
    private String resetCycle;

    public String getName() {
        return name;
    }

    public int getCapacityGb() {
        return capacityGb;
    }

    public double getUsedGb() {
        return usedGb;
    }

    public String getResetCycle() {
        return resetCycle;
    }

    @Override
    public String toString() {
        return "TrafficPackage{" +
                "name='" + name + '\'' +
                ", capacityGb=" + capacityGb +
                ", usedGb=" + usedGb +
                ", resetCycle='" + resetCycle + '\'' +
                '}';
    }
}
