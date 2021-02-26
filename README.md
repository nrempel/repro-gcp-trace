# Reproduction of Google Cloud Trace issue

Deploy this project to Cloud Run.

Enable a Load balancer to point traffic at this Cloud Run deployment.

`curl 'http://<load_balancer_ip>/hello' -H "X-Cloud-Trace-Context: 8379068018630e71f70d3a8fba27724e/1;o=1"`

Change the trace id each time.

Notice the following behaviour:

![](/screenshot1.png)

The internal spans are associated correctly. However, the `load_balancer` and `cloud_run_revision` resources are missing from the flame graph. They **are** present in the logs if you enable logs in the UI:

![](/screenshot2.png)

These two components should show up as their own spans in the flame graph.

You can confirm that the traceid is set correctly by viewing the trace logs:

![](/screenshot3.png)

The logs exist for this trace_id:

![](/screenshot4.png)
