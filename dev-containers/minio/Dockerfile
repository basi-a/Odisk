FROM quay.io/minio/minio:RELEASE.2024-04-18T19-09-19Z

COPY ./cert/server.key /root/.minio/certs/CAs/private.key
COPY ./cert/server.key /root/.minio/certs/private.key
COPY ./cert/server.crt /root/.minio/certs/public.crt