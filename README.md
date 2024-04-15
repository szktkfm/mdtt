# 🗓️ mdtt: Markdown Table Editor TUI

Editing markdown tables can be tedious. `mdtt` simplifies this with a TUI designed for direct terminal edits using vim-like keybindings. It supports outputting to stdout or inplace editing of markdown files.

## 📦 Installation

To install mdtt, just install it with go:

```sh
go install github.com/szktkfm/mdtt@latest
```

Or, download it:

[GitHub Releases](https://github.com/szktkfm/mdtt/releases)

## 🎬 Usage
To start editing your markdown table, simply run:

```sh
mdtt filename.md
```
When launched, mdtt will display the tables from your markdown file in a TUI. 

While editing, you can utilize the following vim-like keybindings to navigate and modify your tables efficiently:

- Navigation: Use `hjkl` to move left, down, up, and right respectively.
- Editing: Enter insert mode by pressing `i` to start editing the cell content, and return to normal mode with `esc` or `ctrl+c`.
- Row and Column Manipulation:
    - Add a new row with `o`.
    - Add a new column with `vo`.
    - Delete the current row or column with `dd`, `vd`.

ここにgif挿入
[]

If you prefer to edit the file directly and save changes back to the same file, use:

```sh
mdtt -i filename.md
```

ここにgif挿入
[]

To pipe contents into mdtt:

```sh
pbpaste | mdtt | pbcopy
```

ここにgif挿入
[]

To create a new table without using an existing file, simply run mdtt without any arguments:

```sh
mdtt
```

ここにgif挿入
[]

If multiple tables are present, you'll be greeted with a view to select which table you want to edit.

```sh
mdtt multiple.md
```

## ⌨️ Key Bindings
The TUI supports the following keybindings for efficient table manipulation:

| Key            | Action            |
| -------------- | ----------------- |
| `↑`/`k`        | Move up           |
| `↓`/`j`        | Move down         |
| `←`/`h`        | Move left         |
| `b`/`pgup`     | Page up           |
| `f`/`pgdn`     | Page down         |
| `ctrl+u`       | Half page up      |
| `ctrl+d`       | Half page down    |
| `g`/`home`     | Go to start       |
| `G`/`end`      | Go to end         |
| `i`            | Insert mode       |
| `esc`/`ctrl+c` | Normal mode       |
| `o`/`v+o`      | Add row/column    |
| `dd`/`v+d`     | Delete row/column |
| `y`            | Copy row          |
| `p`            | Paste             |
| `q`            | Quit              |
| `?`            | Toggle help       |

## 📝 Features
- [x] **Vim-like Keybindings**: Navigate and edit tables using familiar vim commands.
- [x] **Inplace Editing**: Directly modify your original markdown files with the -i option.
- [x] **Piping Support**
- [x] **Multi-Table Selection**
- [ ] **HTML in Cells**: Enable rich content formatting by using HTML directly within table cells.


## 🙏 Acknowledgments
This project, mdtt, was inspired by [mdvtbl](https://github.com/karino2/mdvtbl), a tool that reads markdown from stdin, allows for table editing in a web view, and outputs to stdout. 

## License
[MIT](./LICENSE)
