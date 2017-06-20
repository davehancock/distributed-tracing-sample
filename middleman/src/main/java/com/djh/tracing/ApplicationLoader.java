package com.djh.tracing;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.stream.annotation.EnableBinding;
import org.springframework.cloud.stream.messaging.Source;
import org.springframework.integration.annotation.Gateway;
import org.springframework.integration.annotation.MessagingGateway;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author David Hancock
 */
@SpringBootApplication
public class ApplicationLoader {

    private final Logger LOG = LoggerFactory.getLogger(ApplicationLoader.class);

    public static void main(String[] args) {
        SpringApplication.run(ApplicationLoader.class);
    }

    @RestController
    @EnableBinding(Source.class)
    public class RestEndpointHandler {

        @Autowired
        private PayloadProcessor processor;

        @PutMapping("/foo")
        public void receiveMessage(@RequestBody String message) {
            LOG.info("[Middleman] Received a message to forward on [{}]", message);
            processor.publishPayload(message);
        }
    }

    @MessagingGateway
    public interface PayloadProcessor {

        @Gateway(requestChannel = Source.OUTPUT)
        void publishPayload(String payload);
    }

}
