import { TickerRequest, TickerResponse, CompanyResponse, PriceResponse } from '../protos/watcher_pb.js';
import { WatcherClient } from '../protos/watcher_grpc_web_pb.js';

var watcherService = new WatcherClient('http://'+ window.location.hostname + ':8080',
                               null, null);

var request = new TickerRequest();
request.setTicker('TSLA')

// Unary gRPC call
watcherService.echo(request, {}, (err, response) => {
    if(err) {
        console.log(err.code);
        console.log(err.message);
    } else {
        console.log(response.getTicker());
    }
});
