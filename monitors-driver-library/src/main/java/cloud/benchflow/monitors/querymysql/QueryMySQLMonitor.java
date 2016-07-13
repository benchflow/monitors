package cloud.benchflow.monitors.querymysql;

import cloud.benchflow.monitors.Monitor;
import cloud.benchflow.monitors.MonitorAPI;
import com.google.gson.Gson;
import com.google.gson.annotations.SerializedName;
import com.sun.faban.driver.transport.hc3.ApacheHC3Transport;

import java.util.Map;
import java.util.logging.Logger;

/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public class QueryMySQLMonitor extends Monitor {

    private String completionQuery;
    private String completionQueryValue;
    private String completionQueryMethod;

    public QueryMySQLMonitor(ApacheHC3Transport http,
                             Map<String, String> params,
                             String endpoint,
                             MonitorAPI api) {
        this(http, params,endpoint,api,Logger.getLogger(QueryMySQLMonitor.class.getName()));
    }

    public QueryMySQLMonitor(ApacheHC3Transport http,
                             Map<String, String> params,
                             String endpoint,
                             MonitorAPI api,
                             Logger logger) {
        super(http, params, endpoint, api, logger);
        this.completionQuery = params.get("COMPLETION_QUERY");
        this.completionQueryValue = params.get("COMPLETION_QUERY_VALUE");
        this.completionQueryMethod = params.get("COMPLETION_QUERY_METHOD");
    }

    @Override
    public void monitor() throws Exception {


        String apiUrl = endpoint + "/" + api.getMonitor() + "?"
                        + "query=" + completionQuery + "&value=" + completionQueryValue
                        + "&method=" + completionQueryMethod;

        while(true) {

            String rawResponse = http.fetchURL(apiUrl).toString();
            MonitorQueryResponse response = new Gson().fromJson(rawResponse, MonitorQueryResponse.class);
            if (response.isResult()) {
                break;
            }
            else {
                sleep(1000);
            }
        }

    }

    @Override
    protected void start() throws Exception {
        http.fetchURL(endpoint + "/" + api.getStart());
    }

    @Override
    protected void stop() throws Exception{
        http.fetchURL(endpoint + "/" + api.getStop());
    }

    private static class MonitorQueryResponse {

        @SerializedName("Query_response")
        private int queryResponse;

        @SerializedName("Result")
        private boolean result;

        public MonitorQueryResponse(int queryResponse, boolean result) {
            this.queryResponse = queryResponse;
            this.result = result;
        }

        public int getQueryResponse() {
            return queryResponse;
        }

        public boolean isResult() {
            return result;
        }
    }

}
