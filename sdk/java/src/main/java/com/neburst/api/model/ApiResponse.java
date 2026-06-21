package com.neburst.api.model;

/**
 * Generic API response envelope.
 * All Neburst API responses follow the format: {"code": 0, "msg": "", "data": T}
 *
 * @param <T> the type of the data field
 */
public class ApiResponse<T> {

    private int code;
    private String msg;
    private T data;

    public ApiResponse() {
    }

    public ApiResponse(int code, String msg, T data) {
        this.code = code;
        this.msg = msg;
        this.data = data;
    }

    public int getCode() {
        return code;
    }

    public String getMsg() {
        return msg;
    }

    public T getData() {
        return data;
    }

    public boolean isSuccess() {
        return code == 0;
    }

    @Override
    public String toString() {
        return "ApiResponse{code=" + code + ", msg='" + msg + "', data=" + data + "}";
    }
}
