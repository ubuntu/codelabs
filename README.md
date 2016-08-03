# Our snap/snapcraft codelabs

This is our snap and snapcraft codelabs, fetched from google doc
content.

This is the compiled containing compiled assets or codelabs.
The source branch is at https://github.com/ubuntu/codelabs-source.

## Run the binary assets

Once you are on the codelabs repo, you can just run the simple webserver
from the main repo:

 * ./server
   You can specify the port with -p <port_number>
 * There is snap available name "snap-codelabs" which will run on your localhost,
   port 8123 by default. You can install it with: sudo snap install snap-codelabs
 * If you have polymer-cli (npm install -g polymer-cli), you can just run: polymer serve.

## Add/Update/Remove codelabs

You can use ./codelabs binary which will fetch needed dependencies for you to
add/update or remove codelabs.

 * Adding a new codelabs is as simple as: `./codelabs add <google_doc_id>`.
You can add multiple docs at the same time.
 * Refreshing all codelabs is `./codelabs update`
 * Remove a codelab is `./codelabs remove <google_doc_id|codelab_name>.
You can remove multiple docs at the same time.

You can use -ga <google_analytics> to override the default GA account.

Codelabs are located in `src/codelabs`. All metadata are then regenerated for the website
to pick up.

Do not forget to add/commit and push to the `codelabs` branch each time you
generate or refresh the codelabs assets.

## Tweak category theming and events

The theming is for categories are located in `categories-events.json`.
Only categories available there will be shown in the dropdown filters.

Adding events enables to get events/<event_name> page, which is filtering
codelabs for which one tags match this event.
Images are relative path to images/events/.
