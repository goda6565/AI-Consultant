import axios, {
  AxiosHeaders,
  type AxiosRequestConfig,
  type AxiosRequestHeaders,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import { auth } from "@/shared/config/firebase";
import { env } from "../config";
import { handleApiError } from "./utils";

const adminApiClientInstance = axios.create({
  baseURL: env.NEXT_PUBLIC_ADMIN_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Attach Firebase ID token to each request dynamically
adminApiClientInstance.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const user = auth.currentUser;
    if (user) {
      const token = await user.getIdToken();
      const normalizedHeaders =
        config.headers instanceof AxiosHeaders
          ? config.headers
          : new AxiosHeaders(config.headers as AxiosRequestHeaders);
      normalizedHeaders.set("Authorization", `Bearer ${token}`);
      config.headers = normalizedHeaders;
    }
    return config;
  },
);

export const adminApiClient = async <T>(
  config: AxiosRequestConfig,
): Promise<T> => {
  try {
    const response: AxiosResponse<T> = await adminApiClientInstance(config);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
};

const agentApiClientInstance = axios.create({
  baseURL: env.NEXT_PUBLIC_AGENT_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Attach Firebase ID token to each request dynamically
agentApiClientInstance.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const user = auth.currentUser;
    if (user) {
      const token = await user.getIdToken();
      const normalizedHeaders =
        config.headers instanceof AxiosHeaders
          ? config.headers
          : new AxiosHeaders(config.headers as AxiosRequestHeaders);
      normalizedHeaders.set("Authorization", `Bearer ${token}`);
      config.headers = normalizedHeaders;
    }
    return config;
  },
);

export const agentApiClient = async <T>(
  config: AxiosRequestConfig,
): Promise<T> => {
  try {
    const response: AxiosResponse<T> = await agentApiClientInstance(config);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
};
