param app resource 'radius.dev/Application@v1alpha3'

resource container 'radius.dev/Application/Container@v1alpha3' = {
  name: '${app.name}/container'
  properties: {
    container: {
      image: 'nginx:latest'
    }
  }
}

resource backendhttp 'radius.dev/Application/HttpRoute@v1alpha3' = {
  name: '${app.name}/backendhttp'
}

output test resource 'radius.dev/Application/HttpRoute@v1alpha3' = backendhttp