# Kyma

Create presentations in the terminal from markdown files with fancy transitions.

![slideshow](slideshow.gif)

## Install

```bash
go install github.com/ploMP4/kyma@latest
```

## How it works

Write your slides in plain Markdown. Separate each slide using `----`.

You can customize each slide with optional metadata at the top using YAML front matter:

```yaml
---
transition: swipeLeft
style:
  layout: center
  border: rounded
  border_color: "#880808"
---
```

### Slide Configuration Options

#### Transitions

Control how the presentation moves between slides.

Available transitions:

- `swipeLeft`
- `swipeRight`
- `flip`
- `slideDown`
- `slideUp`

> Leave the transition field empty if you don’t want any transition.

#### Style Options

Customize the look and layout of your slides.

##### Layout

Control alignment both horizontally and vertically:

```yaml
style:
  layout: <horizontal> <vertical>
```

Available options:

- `center`
- `right`
- `left`
- `top`
- `bottom`

> If you provide only one value it will be applied on both horizontal and vertical.

##### Border

Choose a border style for each slide:

- `rounded`
- `double`
- `thick`
- `normal`
- `block`
- `innerHalfBlock`
- `outerHalfBlock`
- `hidden`

##### Border Color

You can use either ANSI color names or HEX codes (e.g. "#ff0000").

```yaml
style:
  border_color: "#880808"
```

## Roadmap

- Add support for more style options like text color and background color
- Allow choosing from any glamour themes
- Create grid-based slide layouts with transitions for each pane
- Add more transition effects
- Support image rendering in terminals (e.g., via the Kitty protocol)

## Contributing

All contributions are welcome.

If you’re planning a significant change or you're unsure about an idea, please open an issue first so we can discuss it in detail.
