import type { paths } from "@/lib/api/openapi.d.ts";
import { GetParameterCommand, SSMClient } from "@aws-sdk/client-ssm";
import createFetchClient, { type Middleware } from "openapi-fetch";
import createClient from "openapi-react-query";

let _initialized = false;
let client: SSMClient | null = null;

const getClient = () => {
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
  const ssmClient = getClient();
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

const myMiddleware: Middleware = {
  async onRequest({ request }) {
    try {
      const apiKey = await getApiKey();
      if (apiKey) {
        request.headers.set("X-API-Key", apiKey);
      }
      return request;
    } catch {
      throw new Error("Unable to authenticate");
    }
  },
};

export const fetchClient = createFetchClient<paths>({
  baseUrl: "http://localhost:8080/api/v1",
});
fetchClient.use(myMiddleware);

export const $api = createClient(fetchClient);
