# mdtt: Markdown Table Editor TUI

Editing markdown tables can often be tedious and error-prone. mdtt is here to change that by providing a Terminal User Interface (TUI) specifically designed for editing markdown tables. By reading markdown files and allowing users to edit them with vim-like keybindings directly in the terminal, mdtt simplifies the process of managing and updating tables. This tool supports both outputting directly to standard output or inplace editing of the original markdown file.

## Installation

To install mdtt, ensure you have Go installed on your system, and then run the following command:

```sh
go install github.com/szktkfm/mdtt@latest
```

## Usage
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

If you prefer to edit the file directly and save changes back to the same file, use:

```sh
mdtt -i filename.md
```

To pipe contents into mdtt:

```sh
cat filename.md | mdtt
```

To create a new table without using an existing file, simply run mdtt without any arguments:

```sh
mdtt
```

If multiple tables are present, you'll be greeted with a view to select which table you want to edit.

```sh
mdtt multiple.md
```


## Key Bindings
The TUI supports the following keybindings for efficient table manipulation:


| Key        | Action            |
| ---------- | ----------------- |
| ↑/k        | Move up           |
| ↓/j        | Move down         |
| ←/h        | Move left         |
| →/l        | Move right        |
| b/pgup     | Page up           |
| f/pgdn     | Page down         |
| ctrl+u     | Half page up      |
| ctrl+d     | Half page down    |
| g/home     | Go to start       |
| G/end      | Go to end         |
| i          | Insert mode       |
| esc/ctrl-c | Normal mode       |
| o/v+o      | Add row/column    |
| dd/v+d     | Delete row/column |
| y          | Copy row          |
| p          | Paste             |
| q          | Quit              |
| ?          | Toggle help       |


## Contributing
Contributions to mdtt are welcome! Feel free to fork the repository, make changes, and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
MIT

