package com.neburst.api.model;

import com.google.gson.annotations.SerializedName;

/**
 * Represents a billing invoice.
 */
public class Invoice {

    private String uuid;
    private double amount;
    private String status;
    private String category;

    @SerializedName("created_at")
    private String createdAt;

    @SerializedName("due_at")
    private String dueAt;

    public String getUuid() {
        return uuid;
    }

    public double getAmount() {
        return amount;
    }

    public String getStatus() {
        return status;
    }

    public String getCategory() {
        return category;
    }

    public String getCreatedAt() {
        return createdAt;
    }

    public String getDueAt() {
        return dueAt;
    }

    @Override
    public String toString() {
        return "Invoice{" +
                "uuid='" + uuid + '\'' +
                ", amount=" + amount +
                ", status='" + status + '\'' +
                ", category='" + category + '\'' +
                ", createdAt='" + createdAt + '\'' +
                ", dueAt='" + dueAt + '\'' +
                '}';
    }
}
