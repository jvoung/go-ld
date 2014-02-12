// Copyright (c) 2014, Jan Voung
// All rights reserved.

// Determine which file resolves the undef symbols of another (required) file.

package main

func ResolveSymbols(f_syms []SymbolTable) []SymLinkInfo {
	imports_exports := make([]SymLinkInfo, 0, len(f_syms))

	// 1. Get the set of defined and undefined syms.
	for _, syms := range f_syms {
		imports_exports = append(imports_exports,
			GetSymLinkInfo(syms))
	}

	// 2. For each undef sym, search through other files
	// to see who defines the same symbol.
	for cur_file, ie := range imports_exports {
		cur_symtab := &f_syms[cur_file]
		for undef_index, _ := range ie.UndefinedSyms {
			sym_name := (*cur_symtab)[undef_index].St_name
			for other_file, other_ie := range imports_exports {
				if other_file == cur_file {
					continue
				}
				def_index, ok := other_ie.ExportedSymHash[sym_name]
				if ok {
					ie.UndefinedSyms[undef_index] = Resolver{
						other_file, def_index}
				}
			}
		}
	}
	return imports_exports
}
