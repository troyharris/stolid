#Stolid - Static Website Suite

This is still in the process of being built. The goal is to be able to have a development directory such as:

```
/
|-- config.json
|-- content
    |-- category_directory
        |-- article.md
    |-- another_category
        |--another_article.md
|-- template
    |-- (template html files)
    |-- css (template css files)
```

You would edit config.json and add path where the generated website should live, the website url, etc.

You could then run `stolid /development/directory` and the site would be generated and an http server would be started. Updating is simple. Make your changes in the development directory and `touch update` in the root directory. Stolid will generate the site from scratch on the fly. 