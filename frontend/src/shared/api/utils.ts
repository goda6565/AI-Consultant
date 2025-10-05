import { AxiosError } from "axios";
import { ErrorCode, type ErrorResponse } from "./admin/model";

export const handleApiError = (error: unknown): ErrorResponse => {
  if (error instanceof AxiosError) {
    const status = error.response?.status;
    const message =
      error.response?.data?.message ||
      error.message ||
      "Unknown error occurred";

    // ステータスコードに基づいてErrorCodeを決定
    let code: ErrorCode;
    switch (status) {
      case 400:
        code = ErrorCode.NUMBER_400;
        break;
      case 401:
        code = ErrorCode.NUMBER_401;
        break;
      case 403:
        code = ErrorCode.NUMBER_403;
        break;
      case 404:
        code = ErrorCode.NUMBER_404;
        break;
      case 409:
        code = ErrorCode.NUMBER_409;
        break;
      case 500:
        code = ErrorCode.NUMBER_500;
        break;
      default:
        code = ErrorCode.NUMBER_500; // デフォルトは500
    }

    return {
      message,
      code,
    };
  }
  throw error;
};
