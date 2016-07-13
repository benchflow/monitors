package cloud.benchflow.monitors;

import com.sun.faban.driver.transport.hc3.ApacheHC3Transport;

import java.util.Map;
import java.util.logging.Logger;


/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public abstract class Monitor {

    protected Map<String, String> params;
    protected String endpoint;
    protected MonitorAPI api;
    protected ApacheHC3Transport http;
    protected final Logger logger;

    public Monitor(ApacheHC3Transport http,
                   Map<String, String> params,
                   String endpoint,
                   MonitorAPI api,
                   Logger fabanLogger) {
        this.endpoint = endpoint;
        this.params = params;
        this.api = api;
        this.http = http;
        this.logger = fabanLogger;
    }

    public abstract void monitor() throws Exception;

    protected abstract void start() throws Exception;

    protected abstract void stop() throws Exception;

    final public void sleep(long millis) {
        try {
            Thread.sleep(millis);
        } catch (InterruptedException e) {
            Thread t = Thread.currentThread();
            t.getUncaughtExceptionHandler().uncaughtException(t, e);
        }
    }

    final public void run() throws Exception {
        start();
        monitor();
        stop();
    }

}
