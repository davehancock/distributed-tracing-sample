package com.djh.tracing;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.sleuth.instrument.web.client.TraceWebClientAutoConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestOperations;
import org.springframework.web.client.RestTemplate;

/**
 * @author David Hancock
 */
@EnableScheduling
@SpringBootApplication
public class ApplicationLoader {

    private final Logger LOG = LoggerFactory.getLogger(ApplicationLoader.class);

    public static void main(String[] args) {
        SpringApplication.run(ApplicationLoader.class);
    }

    /**
     * We explicitly declare this as a Bean to allow Spring to enrich RestTemplate
     * with a TraceRestTemplateInterceptor - {@link TraceWebClientAutoConfiguration}
     */
    @Bean
    public RestOperations restOperations() {
        return new RestTemplate();
    }

    @Component
    public class MessageProducer {

        private RestOperations restOperations;

        public MessageProducer(RestOperations restOperations) {
            this.restOperations = restOperations;
        }

        @Scheduled(initialDelay = 3000, fixedRate = 3000)
        public void produceMessage() {

            String payload = "Payload Contents: " + Math.random();
            LOG.info("[Publisher] Sending message with payload: [{}]", payload);
            restOperations.put("http://middleman:8080/foo", payload);
        }
    }

}
