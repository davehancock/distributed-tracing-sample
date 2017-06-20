package com.djh.tracing;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.stream.annotation.EnableBinding;
import org.springframework.cloud.stream.messaging.Sink;
import org.springframework.integration.annotation.ServiceActivator;

/**
 * @author David Hancock
 */
@SpringBootApplication
public class ApplicationLoader {

    private static final Logger LOG = LoggerFactory.getLogger(ApplicationLoader.class);

    public static void main(String[] args) {
        SpringApplication.run(ApplicationLoader.class);
    }

    @EnableBinding(Sink.class)
    public class MessageSink {

        @ServiceActivator(inputChannel = Sink.INPUT)
        public void consumeMessage(Object payload) {
            LOG.info("[Consumer] Received payload: [{}] ", payload);
        }
    }

}
