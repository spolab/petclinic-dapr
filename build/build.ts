import Client, { connect } from "@dagger.io/dagger";

connect(async(client: Client) => {
    console.info(`Successfully connected to the client ${client.clientHost}`);
    const golang = client.container().from("golang:1.19").exec(["go", "--version"]);
    let output = await golang.stdout().contents();
    console.info(await golang.exitCode());
    console.info(output);
});