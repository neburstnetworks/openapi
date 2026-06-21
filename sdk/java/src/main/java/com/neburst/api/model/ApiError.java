package com.neburst.api.model;

/**
 * Represents an API error returned by the Neburst API.
 * Thrown as a {@link com.neburst.api.NeburstApiException} when the response code is non-zero.
 */
public class ApiError {

    private final int code;
    private final String message;

    public ApiError(int code, String message) {
        this.code = code;
        this.message = message;
    }

    public int getCode() {
        return code;
    }

    public String getMessage() {
        return message;
    }

    @Override
    public String toString() {
        return "ApiError{code=" + code + ", message='" + message + "'}";
    }
}
