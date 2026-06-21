package com.neburst.api.model;

/**
 * Represents user account balance.
 */
public class Balance {

    private double available;
    private double locked;
    private String currency;

    public double getAvailable() {
        return available;
    }

    public double getLocked() {
        return locked;
    }

    public String getCurrency() {
        return currency;
    }

    @Override
    public String toString() {
        return "Balance{" +
                "available=" + available +
                ", locked=" + locked +
                ", currency='" + currency + '\'' +
                '}';
    }
}
