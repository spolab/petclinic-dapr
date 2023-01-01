/*
Copyright 2022 Alessandro Santini

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

// The gateway is responsible for turning a REST request into a command and forward it to the relevant actor.
// The actor will respond with either a) a series of events to be broadcasted or b) an error indicating if the command was invalid or an error occurred while executing it.
// If the actor execution is successful, the gateway will be responsible for broadcasting the event to all the parties interested.
func main() {

}
