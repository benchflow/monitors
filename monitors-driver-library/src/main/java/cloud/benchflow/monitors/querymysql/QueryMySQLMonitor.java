package cloud.benchflow.monitors.querymysql;

import cloud.benchflow.monitors.Monitor;
import cloud.benchflow.monitors.MonitorAPI;
import com.google.gson.Gson;
import com.google.gson.annotations.SerializedName;

import org.apache.http.client.fluent.Request;

import java.nio.charset.Charset;
import java.util.Map;
import java.util.logging.Logger;

/**
 * @author Simone D'Avico (simonedavico@gmail.com)
 *
 * Created on 13/07/16.
 */
public class QueryMySQLMonitor extends Monitor {

    protected final static int SQLMONITOR_SLEEP_TIME = 5000;

    private String completionQuery;
    private String completionQueryValue;
    private String completionQueryMethod;

    public QueryMySQLMonitor(Map<String, String> params,
                             String endpoint,
                             MonitorAPI api,
                             Logger logger) {
        super(params, endpoint, api, logger);
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
            String rawResponse = Request.Get(apiUrl).execute().returnContent().asString(Charset.forName("UTF-8"));
            MonitorQueryResponse response = new Gson().fromJson(rawResponse, MonitorQueryResponse.class);
            if (response.isResult()) {
                break;
            }
            else {
                //TODO: we could use the query response to regulate the sleep amount
                sleep(SQLMONITOR_SLEEP_TIME);
            }
        }

    }

    @Override
    public void start() throws Exception {
        Request.Get(endpoint + "/" + api.getStart());
    }

    @Override
    public void stop() throws Exception{
        Request.Get(endpoint + "/" + api.getStop());
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
