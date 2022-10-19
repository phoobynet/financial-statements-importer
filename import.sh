#!/bin/zsh

sqlite3 2022q3 <<EOF
.mode tabs
.import sub.txt subs
.import tag.txt tags
.import pre.txt pre
.import num.txt nums
create index if not exists idx_subs_cik on subs (cik);
create index if not exists idx_subs_adsh on subs (adsh);
create index if not exists idx_nums_adsh on nums (adsh);
create index if not exists idx_pre_adsh on pre (adsh);
EOF
