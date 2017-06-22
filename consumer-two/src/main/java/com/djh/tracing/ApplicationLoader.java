package com.djh.tracing;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author David Hancock
 */
@SpringBootApplication
public class ApplicationLoader {

    private static final Logger LOG = LoggerFactory.getLogger(ApplicationLoader.class);

    public static void main(String[] args) {
        SpringApplication.run(ApplicationLoader.class);
    }

    @RestController
    public class RestEndpointHandler {

        @GetMapping("/bar")
        public String receiveMessage() {

            // Introduce some random failure here...
            if (Math.random() > 0.5) {
                LOG.info("[consumer-two] Freaking Out!");
                throw new RuntimeException("Uncontrollable cascading exception");
            }

            LOG.info("[consumer-two] Approving Request");
            return "[Approved]";
        }
    }

}
