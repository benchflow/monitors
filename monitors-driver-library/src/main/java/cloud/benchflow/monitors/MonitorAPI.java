package cloud.benchflow.monitors;

/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public final class MonitorAPI {

    private String start;
    private String monitor;
    private String stop;

    public MonitorAPI(String start, String monitor, String stop) {
        this.start = start;
        this.monitor = monitor;
        this.stop = stop;
    }

    public String getStart() {
        return start;
    }

    public String getMonitor() {
        return monitor;
    }

    public String getStop() {
        return stop;
    }
}
