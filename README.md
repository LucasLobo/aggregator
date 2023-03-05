# Aggregator CLI

Simple command line application that calculates an aggregated metric based on the input events. It currently supports 
calculating the moving average for the last X minutes.

The input events can be fetched from a file or from AWS SQS. The output events can be written to a file or to the stdout.

## Context

One of the metrics we use for our clients' SLAs is the delivery time of a translation. In the context of this problem,
our translation flow is modeled as only one event.

### *translation_delivered*

Example:

```json
{
    "timestamp": "2018-12-26 18:12:19.903159",
    "translation_id": "5aa5b2f39f7254a75aa4",
    "source_language": "en",
    "target_language": "fr",
    "client_name": "airliberty",
    "event_name": "translation_delivered",
    "duration": 20,
    "nr_words": 100
}
```

## Quick Start

To calculate, for each minute, the moving average delivery time of all translations for the past 10 minutes you should 
call the application like this:

	./aggregator-cli moving-average --window_size 10 --input_file example/events.json --output_folder example/output
	
The input file must have the following format:

```
{"timestamp": "2018-12-26 18:11:08.509654","translation_id": "5aa5b2f39f7254a75aa5","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 20}
{"timestamp": "2018-12-26 18:15:19.903159","translation_id": "5aa5b2f39f7254a75aa4","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 31}
{"timestamp": "2018-12-26 18:23:19.903159","translation_id": "5aa5b2f39f7254a75bb3","source_language": "en","target_language": "fr","client_name": "taxi-eats","event_name": "translation_delivered","nr_words": 100, "duration": 54}
```

The lines in the input must be ordered by the `timestamp` key, from lower (oldest) to higher values, just like in the 
example input above.

The output file will have the following format.

```
{"date": "2018-12-26 18:11:00", "average_delivery_time": 0}
{"date": "2018-12-26 18:12:00", "average_delivery_time": 20}
{"date": "2018-12-26 18:13:00", "average_delivery_time": 20}
{"date": "2018-12-26 18:14:00", "average_delivery_time": 20}
{"date": "2018-12-26 18:15:00", "average_delivery_time": 20}
{"date": "2018-12-26 18:16:00", "average_delivery_time": 25.5}
{"date": "2018-12-26 18:17:00", "average_delivery_time": 25.5}
{"date": "2018-12-26 18:18:00", "average_delivery_time": 25.5}
{"date": "2018-12-26 18:19:00", "average_delivery_time": 25.5}
{"date": "2018-12-26 18:20:00", "average_delivery_time": 25.5}
{"date": "2018-12-26 18:21:00", "average_delivery_time": 25.5}
{"date": "2018-12-26 18:22:00", "average_delivery_time": 31}
{"date": "2018-12-26 18:23:00", "average_delivery_time": 31}
{"date": "2018-12-26 18:24:00", "average_delivery_time": 42.5}
```

## Flags

Below are the flags that can be used to configure the tool:

| Flag          | Usage                                                | Mandatory | Note                                                      |
|---------------|------------------------------------------------------|-----------|-----------------------------------------------------------|
| window_size   | Window size to use in the moving average calculation | `true`    | Defaults to 10 if < 1                                     | 
| input_file    | File name where the input events are stored          | `false`   | Either `input_file` or `queue_url` must be provided       | 
| queue_url     | SQS Queue from which to read the events              | `false`   | Either `input_file` or `queue_url` must be provided       | 
| output_folder | Folder where output events will be written to        | `false`   | If none is provided, output will be printed to the stdout | 

## Reading from AQS SQS Queue
To read from an AWS SQS queue you must:

1. Create an AWS SQS Queue (FIFO).
2. Authenticate AWS locally. Follow [these instructions from AWS Docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html). `LoadDefaultConfig` is used to get the values.
3. Run the CLI like this `./aggregator-cli moving-average --window_size 10 --queue_url QUEUE_URL --output_folder example/output`
4. Add messages to queue. Each message should have the same format as one of the input lines.

## Example Input

An example input file is provided in `example/input.json`.

## Project Structure

This project is *slightly* (wink) overengineered for this purpose. The motivation for this was not only to complete 
the code challenge but also to experiment with the hexagonal architecture pattern.

The hexagonal architecture is used isolate the business logic from the input and output. The project can be easily 
extended to read from other input sources and to emit events with the results instead of writing them to files/stdout 
(e.g., using AWS SNS). It's also easy to create new aggregation methods besides moving average.

```
┌── cmd                           // Commands to run the application
└── internal                      // Everything that is internal to the application
    ├── common                    // Common modules used in multiple places within application
    ├── core                      // Core (business logic)
        ├── application           // Business logic. Implements interactorprt and uses infrastructureprt
        ├── domain                // Business models
        ├── infrastructureprt     // Ports (interfaces) for infrastructure
        └── interactorprt         // Ports (interfaces) for interactors
    ├── infrastructure            // Primary/Driving adapters - entrypoint (e.g., files, queues, api, etc)
    └── interactors               // Secondary/Driven adapters (e.g., files, databases, etc). 
```

## FAQ

**Q1: Is so much code really needed?**  
A1: The answer to this question is a simple no. I wanted to take this opportunity to implement a new project
from scratch using a Hexagonal Architecture, and while it is definitely overengineered for this purpose, it also allows
me to experiment and showcase my architectural skills.

**Q2: How does the algorithm for calculating the moving averages work?**  
A2: The algorithm uses the sliding window technique that traverses each minute-time-bucket between the first and last
events. It has a running total/count, which it uses to calculate the current average at each step. The head and tail are
used to add/remove from the running values. This allows us to process each time bucket only once. At the bottom of this
page you can find a diagram that explains how the algorithm works.

**Q3: What is the time and space complexity of this algorithm?**  
A3: This algorithm has a time complexity of O(N), where N is the number of minutes between the first and last event.
This is because we traverse each time-bucket between the first and last event only once (e.g., if there are 2 events
spaced 1 hour apart, we will need around 60 iterations). The space complexity is O(K), where K is the window size,
since we store the partial averages for each time-bucket in the window.

**Q4: How would you improve the algorithm?**  
A4.1: With this algorithm we only calculate the moving average for each time-bucket the moment an event arrives. This
means that in a real-world scenario we would be idle for that time. We could pre-configure a maximum threshold of
waiting time to starting processing the next time-bucket (e.g., if we have two events one hour apart, we would be
waiting for 60 minutes and then process 60 time-buckets at once. If we add a threshold of 3 minutes, we could process 
time-bucket of minute 1 at minute 4, minute 2 at minute 5, etc., instead of waiting another 60 minutes).  
A4.2: Additionally, since we read from the file one-by-one and write to the file one-by-one, we lose some time for each
fetch and each store. We could put the fetch and store processes in separate go-routines and run them concurrently to 
save some IO time.

**Q5: What could you have done to make this even better?**  
A5: There are a couple of things that could have been done to improve this solution, which I thought about but
did not due to a lack of time:
- Emit output events to AWS SNS or similar.
- Adding e2e/integration tests to ensure that the file reading and writing work as well as the algorithm
itself.
- Adding more unit tests. Including in other parts of the architecture (outside the business logic).
- Adding Docker.
- Adding a Makefile.
- Testing the business logic with proper mocks.
- Maybe un-overengineer this if the goal is to use and extend the tool in the future only as a CLI.

![Sliding Window](sliding_window.jpg)