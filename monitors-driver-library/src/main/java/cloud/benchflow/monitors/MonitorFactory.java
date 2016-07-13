package cloud.benchflow.monitors;

import cloud.benchflow.monitors.cpu.CpuMonitor;
import cloud.benchflow.monitors.querymysql.QueryMySQLMonitor;

import com.sun.faban.driver.transport.hc3.ApacheHC3Transport;

import java.util.Map;
import java.util.logging.Logger;

/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public class MonitorFactory {

    public static Monitor getMonitor(ApacheHC3Transport http,
                                     String monitorName,
                                     Map<String, String> parameters,
                                     String endpoint,
                                     String startApi,
                                     String stopApi,
                                     String monitorApi) {

        MonitorAPI api = new MonitorAPI(startApi, monitorApi, stopApi);
        switch (monitorName) {

            case "sqlquery": {
                return new QueryMySQLMonitor(http, parameters, endpoint, api);
            }

            case "cpu": {
                return new CpuMonitor(http, parameters, endpoint, api);
            }

            default: {
                throw new RuntimeException("Cannot find monitor implementation named " + monitorName);
            }
        }
    }

    public static Monitor getMonitor(ApacheHC3Transport http,
                                     String monitorName,
                                     Map<String, String> parameters,
                                     String endpoint,
                                     String startApi,
                                     String stopApi,
                                     String monitorApi,
                                     Logger fabanLogger) {

        MonitorAPI api = new MonitorAPI(startApi, monitorApi, stopApi);
        switch (monitorName) {

            case "sqlquery": {
                return new QueryMySQLMonitor(http, parameters, endpoint, api, fabanLogger);
            }

            case "cpu": {
                return new CpuMonitor(http, parameters, endpoint, api, fabanLogger);
            }

            default: {
                throw new RuntimeException("Cannot find monitor implementation named " + monitorName);
            }
        }
    }

}