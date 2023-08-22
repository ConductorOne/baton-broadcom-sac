FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-broadcom-sac"]
COPY baton-broadcom-sac /