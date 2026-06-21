package com.neburst.api;

import com.neburst.api.model.ApiError;

/**
 * Checked exception thrown when the Neburst API returns a non-zero error code.
 */
public class NeburstApiException extends Exception {

    private final int code;

    public NeburstApiException(int code, String message) {
        super(message);
        this.code = code;
    }

    public NeburstApiException(ApiError error) {
        super(error.getMessage());
        this.code = error.getCode();
    }

    /**
     * Returns the API error code.
     */
    public int getCode() {
        return code;
    }

    @Override
    public String toString() {
        return "NeburstApiException{code=" + code + ", message='" + getMessage() + "'}";
    }
}
