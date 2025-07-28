import { GetParameterCommand, SSMClient } from "@aws-sdk/client-ssm";
import axios from "axios";

let _initialized = false;
let client: SSMClient | null = null;

const getSSMClient = () => {
  if (!_initialized) {
    client = new SSMClient({
      region: "us-east-1",
      endpoint: "http://localhost:4566",
      credentials: {
        accessKeyId: "test",
        secretAccessKey: "test",
      },
    });
    _initialized = true;
  }

  return client;
};

const getApiKey = async () => {
  const ssmClient = getSSMClient();
  if (!ssmClient) throw new Error("missing ssm client");
  const input = { Name: "api_key", WithDecryption: true };
  const command = new GetParameterCommand(input);
  try {
    const res = await ssmClient.send(command);
    return res.Parameter?.Value;
  } catch (err: unknown) {
    console.error(err);
    throw err;
  }
};

export const $api = axios.create({
  baseURL: "/api/v1",
  timeout: 1000 * 30,
});

$api.interceptors.request.use(
  async (config) => {
    const token = await getApiKey();
    if (token) {
      config.headers.Authorization = token;
    }
    return config;
  },
  (error) => {
    console.error("Request error:", error);
    return Promise.reject(error);
  }
);
