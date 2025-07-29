import { useAuth0 } from "@auth0/auth0-react";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { getJob, getJobs } from "./queries";

export const useAuthToken = () => {
  const [token, setToken] = useState("");
  const { getAccessTokenSilently, getAccessTokenWithPopup } = useAuth0();

  useEffect(() => {
    if (!token) {
      (async () => {
        try {
          const t = await getAccessTokenSilently({
            authorizationParams: {
              audience: import.meta.env.VITE_AUTH0_AUDIENCE,
            },
          });
          setToken(t);
        } catch (e) {
          console.error(e);
          try {
            const t = await getAccessTokenWithPopup({
              authorizationParams: {
                audience: import.meta.env.VITE_AUTH0_AUDIENCE,
              },
            });
            if (t) {
              setToken(t);
            }
          } catch (err) {
            console.error(e);
          }
        }
      })();
    }
  }, [getAccessTokenSilently, getAccessTokenWithPopup, token]);

  return token;
};

export const useJobs = (userId?: string) => {
  const token = useAuthToken();
  return useQuery({
    queryKey: ["jobs", userId],
    queryFn: () => getJobs({ token, userId: userId! }),
    enabled: !!userId,
  });
};

export const useJob = (jobId?: string) => {
  const token = useAuthToken();
  return useQuery({
    queryKey: ["job", jobId],
    queryFn: () => getJob({ token, id: jobId! }),
    enabled: !!jobId,
  });
};
