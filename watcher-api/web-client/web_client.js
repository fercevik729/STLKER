import { TickerRequest, TickerResponse, CompanyResponse, PriceResponse } from '../protos/watcher_pb.js';
import { WatcherClient } from '../protos/watcher_grpc_web_pb.js';

var watcherService = new WatcherClient('http://localhost:8080');

var request = new TickerRequest();
request.setTicker('SPY')

watcherService.getInfo(request, {}, (err, response) => {
    if(err) {
        console.log(err.code);
        console.log(err.message);
    } else {
        console.log(response.getPrice());
    }
});