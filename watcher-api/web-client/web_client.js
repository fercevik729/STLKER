const {TickerRequest, TickerResponse, CompanyResponse, PriceResponse} = require('../protos/watcher_pb.js');
const {WatcherClient} = require('../protos/watcher_grpc_web_pb.js');

var watcherService = new WatcherClient('http://localhost:8080');

var request = new TickerRequest();
request.setTicker('SPY')

watcherService.getInfo(request, {}, function(err, response) {
    if(err) {
        console.log(err.code);
        console.log(err.message);
    } else {
        console.log(response.getPrice());
    }
});