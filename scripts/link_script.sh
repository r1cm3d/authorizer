#!/usr/bin/sh
# link_script generates the executable file 'authorizer'
cat <<'EOF' >authorizer
	#!/usr/bin/sh
	docker run --interactive \
		authorizer
EOF