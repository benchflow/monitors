package cloud.benchflow.monitors.querymysql;

import cloud.benchflow.monitors.Monitor;
import cloud.benchflow.monitors.MonitorAPI;
import com.google.gson.Gson;
import com.google.gson.annotations.SerializedName;

import org.apache.http.client.fluent.Request;

import java.net.URI;
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

        apiUrl = new URI(apiUrl).normalize().toString();

        while(true) {

            logger.info("[QueryMySQL] About to query: " + apiUrl);
            
            String rawResponse = Request.Get(apiUrl).execute().returnContent().asString(Charset.forName("UTF-8"));

            logger.info("[QueryMySQL] Raw response: " + rawResponse);

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
        String startApi = new URI(endpoint + "/" + api.getStart()).normalize().toString();
        logger.info("[QueryMySQL] About to start: " + startApi);
        Request.Post(startApi);
    }

    @Override
    public void stop() throws Exception{
        String stopApi = new URI(endpoint + "/" + api.getStop()).normalize().toString();
        logger.info("[QueryMySQL] About to stop: " + stopApi);
        Request.Delete(stopApi);
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
