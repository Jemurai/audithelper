FROM python:3.9

RUN apt-get update && \
    apt-get -y dist-upgrade && \
    apt-get -y install \
    awscli \
    jq \
    && \
    apt -y autoremove && \
    apt clean && \
    useradd -m jemurai

RUN     git clone https://github.com/toniblyx/prowler /usr/local/prowler && \
    pip install \
    aws-shell \
    boto3 \
    cartography \
    cloudsplaining \
    detect-secrets \
    scout

COPY ./profile.d/* /etc/profile.d/
COPY ./scripts/* /usr/local/bin/

USER jemurai
WORKDIR /work
VOLUME /work

CMD [ "/bin/bash", "-l" ]
