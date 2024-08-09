# Use the Ubuntu 22.04 image as a base
FROM ubuntu:22.04

# Set environment variables for versions
ENV SELENIUM_VERSION="4.22.0"
ENV CHROME_DRIVER_VERSION="127.0.6533.99"
ENV CHROME_VERSION="latest"

# Update the package list and install necessary packages
RUN apt-get update && \
    apt-get install -y wget unzip curl xvfb libxi6 libgconf-2-4 openjdk-11-jre

# Download and install Google Chrome
RUN wget --no-verbose https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    dpkg -i google-chrome-stable_current_amd64.deb || apt-get install -y -f

# Clean up
RUN rm google-chrome-stable_current_amd64.deb

# Fetch the latest ChromeDriver version and install it
RUN CHROME_DRIVER_VERSION=$(wget -qO- https://chromedriver.storage.googleapis.com/LATEST_RELEASE) && \
    wget --no-verbose https://chromedriver.storage.googleapis.com/${CHROME_DRIVER_VERSION}/chromedriver_linux64.zip && \
    unzip chromedriver_linux64.zip -d /usr/local/bin && \
    rm chromedriver_linux64.zip && \
    chmod +x /usr/local/bin/chromedriver

# Download and install Selenium Standalone Server
RUN wget --no-verbose https://github.com/SeleniumHQ/selenium/releases/download/selenium-${SELENIUM_VERSION}/selenium-server-${SELENIUM_VERSION}.jar && \
    mv selenium-server-${SELENIUM_VERSION}.jar /usr/local/bin/selenium-server-standalone.jar

# Expose port 4444 for Selenium Server
EXPOSE 4444

# Run Selenium Standalone Server with ChromeDriver and Chrome
CMD xvfb-run -a java -jar /usr/local/bin/selenium-server-standalone.jar standalone --port 4444
