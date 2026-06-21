package com.neburst.api.model;

import java.util.List;

/**
 * Represents traffic information for a compute instance.
 */
public class Traffic {

    private List<TrafficPackage> packages;

    public List<TrafficPackage> getPackages() {
        return packages;
    }

    @Override
    public String toString() {
        return "Traffic{packages=" + packages + '}';
    }
}
