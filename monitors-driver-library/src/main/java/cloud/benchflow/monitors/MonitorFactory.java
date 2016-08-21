package cloud.benchflow.monitors;

import cloud.benchflow.monitors.cpu.CpuMonitor;
import cloud.benchflow.monitors.querymysql.QueryMySQLMonitor;

import java.util.Map;
import java.util.logging.Logger;

/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public class MonitorFactory {

    public static Monitor getMonitor(String monitorName,
                                     Map<String, String> parameters,
                                     String endpoint,
                                     String startApi,
                                     String stopApi,
                                     String monitorApi,
                                     Logger fabanLogger) {

        MonitorAPI api = new MonitorAPI(startApi, monitorApi, stopApi);
        switch (monitorName) {

            case "querymysql": {
                return new QueryMySQLMonitor(parameters, endpoint, api, fabanLogger);
            }

            case "cpu": {
                return new CpuMonitor(parameters, endpoint, api, fabanLogger);
            }

            default: {
                throw new RuntimeException("Cannot find monitor implementation named " + monitorName);
            }
        }
    }

}