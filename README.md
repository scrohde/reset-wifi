# reset-wifi
reset-wifi continuously checks if the connection to the internet is down and
resets the network if necessary.

This works around an issue I'm having with my laptop.

    Usage of reset-wifi:
    	-dst string
    		destination to ping (default "google.com")
    	-monitor duration
    		(duration) delayed start to monitoring (default 1m0s)
    	-ping duration
    		(duration) interval between pings (default 2s)
    	-pwd
    		ask for root password

Duration flags are a sequence of positive decimal numbers, each with optional
fraction and a unit suffix,such as "300ms", "1.5h" or "2h45m". Valid time units
are "ns", "us" (or "µs"), "ms", "s", "m", "h".
