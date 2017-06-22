package com.djh.tracing;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.sleuth.instrument.web.client.TraceWebClientAutoConfiguration;
import org.springframework.cloud.stream.annotation.EnableBinding;
import org.springframework.cloud.stream.messaging.Source;
import org.springframework.context.annotation.Bean;
import org.springframework.integration.annotation.Gateway;
import org.springframework.integration.annotation.MessagingGateway;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestOperations;
import org.springframework.web.client.RestTemplate;

/**
 * @author David Hancock
 */
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

    @RestController
    @EnableBinding(Source.class)
    public class RestEndpointHandler {

        private MessageProcessor messageProcessor;
        private RestOperations restOperations;

        public RestEndpointHandler(MessageProcessor messageProcessor, RestOperations restOperations) {
            this.messageProcessor = messageProcessor;
            this.restOperations = restOperations;
        }

        @PutMapping("/foo")
        public void receiveMessage(@RequestBody String message) {

            LOG.info("[Middleman] Received a message to forward on [{}]", message);

            // Fire off some async event for some downstream service to consume
            messageProcessor.publishPayload(message);

            // Send a blocking GET request to a service and forward the response to another service
            String anotherPayload = restOperations.getForObject("http://consumer-two:8080/bar", String.class);
            restOperations.put("http://consumer-three:8080/gofoo", anotherPayload);
        }
    }

    @MessagingGateway
    public interface MessageProcessor {

        @Gateway(requestChannel = Source.OUTPUT)
        void publishPayload(String payload);
    }

}
