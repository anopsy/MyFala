# NRPSapp

The purpose of this app is to help the user to find the nearest spots
where the weather and surfing conditions are good,
based on the user's location and specified range/distance from the location.

The app should be able to establish the user's location // TODO
or ask for the start location if the user wants to check surf conditions elsewhere

The app should be able to create a database of surf spots in Europe and around the world //TODO
It should be able to localize the nearest spots//TODO
and check if the surf and weather conditions there are good.//TODO

Good surf conditions will be fixed values at the beginning
primary swell > 0.6 m 
wave period > 5s
wind < 30 km/h
But function to specify the user's level 
or user's preferences should be added later.//TODO

NRSPapp
-cmd
    -server
	      main.go
-pkg
    -api
	      user.go
        spots.go
    -app
    	-server.go
	    -handlers.go
	    -routes.go
    -repository
      -spotlist.go


The app takes location and distance from the location		 // user.go or spots.go
Creates list of the spots in that range				// spotlist.go spots.go
Checks surf & weather condition at each spot 		       // no idea how to do that and where to put that
It checks them using information found on surf&weather forecast services
windy.com
magicseaweed.com 
surf-forecast.com
The app displays the results on a map, as a green dot	//no idea how to do that and where to put that
and  shows detailed weather and surf conditions



//TODO how to download information from forecast services and how to handle them in the app, where to place it in the app structure
//TODO how to display map and the green dot, where to place it in the app structure



