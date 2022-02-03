# ui

# Local Dev

```bash
nerdctl run -ti --rm -v $(pwd):/app --net=host --env NODE_OPTIONS="--openssl-legacy-provider" node:alpine bash
```

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Run your tests
```
npm run test
```

### Lints and fixes files
```
npm run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
