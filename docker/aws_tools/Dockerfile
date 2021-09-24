FROM python:3.9

RUN apt-get update && \
    apt-get -y dist-upgrade && \
    apt-get -y install \
    awscli \
    jq \
    less \
    vim \
    zip \
    && \
    apt -y autoremove && \
    apt clean && \
    useradd -m jemurai

RUN git clone https://github.com/toniblyx/prowler /usr/local/prowler && \
    pip install \
    aws-shell \
    boto3 \
    cartography \
    cloudsplaining \
    detect-secrets \
    scoutsuite && \
    /bin/sh -c "$(curl -fsSL https://raw.githubusercontent.com/turbot/steampipe/main/install.sh)"

COPY ./profile.d/* /etc/profile.d/
COPY ./scripts/* /usr/local/bin/

USER jemurai

RUN  steampipe plugin install aws

WORKDIR /work
VOLUME /work

EXPOSE 9193

CMD [ "/bin/bash", "-l" ]
