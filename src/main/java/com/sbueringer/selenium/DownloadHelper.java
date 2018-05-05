package com.sbueringer.selenium;

import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.RemoteWebDriver;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.UnsupportedEncodingException;
import java.net.URL;
import java.net.URLDecoder;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.HashMap;
import java.util.List;
import java.util.stream.Collectors;

public class DownloadHelper {

    private RemoteWebDriver driver;
    private boolean isDriverInitialied = false;

    public static void main(String[] args) {
        DownloadHelper dl = new DownloadHelper();
        dl.initDriver();
        dl.getSlides("https://kccnceu18.sched.com", "slides/2018-kubecon-eu");
        dl.getSlides("https://kccncna17.sched.com/", "slides/2017-kubecon-na");
        dl.getSlides("https://cloudnativeeu2017.sched.com/", "slides/2017-kubecon-eu");
        dl.disposeRemoteWebDriver();
    }

    private void getSlides(String url, String outputPath) {
        driver.get(url);

        // Find all Sessions
        List<WebElement> sessionLinkWebElements = driver.findElementsByXPath("//span[contains(@class,'event')]/a");
        List<String> sessionLinks = sessionLinkWebElements.stream().map(a -> a.getAttribute("href")).collect(Collectors.toList());

        HashMap<String, List<String>> downloadLinksPerSession = new HashMap<>();

        // Find all attached files
        sessionLinks.forEach(sessionLink -> {
                    driver.get(sessionLink);
                    List<WebElement> sessionWebElement = driver.findElementsByXPath("//span[contains(@class,'event')]");
                    String sessionName = sessionWebElement.get(0).getText();

                    List<WebElement> attachedFiles = driver.findElementsByXPath("//a[contains(@class,'file-uploaded')]");
                    List<String> attachedFilesLinks = attachedFiles.stream().map(attachedFile -> attachedFile.getAttribute("href")).collect(Collectors.toList());

                    if (attachedFilesLinks.size() > 0) {
                        System.out.println("Found " + attachedFilesLinks.size() + " files for session " + sessionName);
                        downloadLinksPerSession.put(sessionName, attachedFilesLinks);
                    }
                }
        );

        new File(outputPath).mkdirs();

        // Download all attached files
        downloadLinksPerSession.forEach((sessionName, downloadLinks) -> downloadLinks.forEach(downloadLink -> {
                    try {
                        String fileName = sessionName + " - " + URLDecoder.decode(downloadLink.substring(downloadLink.lastIndexOf("/") + 1), "UTF-8");
                        fileName = fileName.replace("/","");
                        Path outputFile = Paths.get(outputPath + File.separator + fileName);

                        if (Files.exists(outputFile)) {
                            System.out.println("Skipping " + downloadLink + " because it already exists.");
                        } else {
                            System.out.println("Downloading " + downloadLink + " to " + outputFile);
                            try (InputStream in = new URL(downloadLink).openStream()) {
                                Files.copy(in, outputFile, StandardCopyOption.REPLACE_EXISTING);
                            } catch (IOException e) {
                                e.printStackTrace();
                            }
                        }
                    } catch (UnsupportedEncodingException e) {
                        e.printStackTrace();
                    }
                }
        ));
    }

    private void initDriver() {
        if (!isDriverInitialied) {
            isDriverInitialied = true;
            createRemoteWebDriver();
        } else {
            while (driver == null) {
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }
    }

    private void createRemoteWebDriver() {
        ChromeOptions options = new ChromeOptions();
        driver = new ChromeDriver(options);
    }

    private void disposeRemoteWebDriver() {
        if (driver != null) {
            driver.quit();
        }
    }
}