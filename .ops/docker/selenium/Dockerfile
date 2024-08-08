FROM openjdk:11-jre-slim

RUN apt-get update && apt-get install -y \
    wget \
    unzip \
    curl \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/selenium && \
    wget -q -O /opt/selenium/selenium-server-standalone.jar https://selenium-release.storage.googleapis.com/3.141/selenium-server-standalone-3.141.59.jar

RUN CHROMEDRIVER_VERSION=`curl -sS chromedriver.storage.googleapis.com/LATEST_RELEASE` && \
    wget -q -O /tmp/chromedriver_linux64.zip http://chromedriver.storage.googleapis.com/$CHROMEDRIVER_VERSION/chromedriver_linux64.zip && \
    unzip /tmp/chromedriver_linux64.zip -d /usr/local/bin/ && \
    rm /tmp/chromedriver_linux64.zip

RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y ./google-chrome-stable_current_amd64.deb && \
    rm google-chrome-stable_current_amd64.deb && \
    apt-get -f install -y && \
    rm -rf /var/lib/apt/lists/*

COPY .ops/sh/start_selenium.sh /opt/selenium/start_selenium.sh
RUN chmod +x /opt/selenium/start_selenium.sh

EXPOSE 4444

ENTRYPOINT ["/opt/selenium/start_selenium.sh"]
