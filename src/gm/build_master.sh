#! /bin/bash 

get_script_dir()
{
    SOURCE="$0"
    while [ -h "$SOURCE"  ]; do # resolve $SOURCE until the file is no longer a symlink
        DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
        SOURCE="$(readlink "$SOURCE")"
        [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
    done
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    echo $DIR
}

script_dir=`get_script_dir`
echo "script dir ${script_dir}"

go build -o gm

mkdir -p /etc/gm

rm -rf /etc/gm/conf
rm -rf /etc/gm/scripts
rm -rf /usr/bin/gm
rm -rf /usr/sbin/gm

mkdir -p /etc/gm/conf

ln -s $script_dir/gm /usr/bin/gm
if [ $? == 0 ]; then 
    echo "ln -s $script_dir/gm /usr/bin/gm succeed"
else
    echo "ln -s $script_dir/gm /usr/bin/gm failed"
fi

ln -s $script_dir/gm /usr/sbin/gm
if [ $? == 0 ]; then 
    echo "ln -s $script_dir/gm /usr/sbin/gm succeed"
else
    echo "ln -s $script_dir/gm /usr/sbin/gm failed"
fi

ln -s $script_dir/../../scripts /etc/gm/scripts
if [ $? == 0 ]; then 
    echo "ln -s $script_dir/../../scripts /etc/gm/scripts succeed"
else
    echo "ln -s $script_dir/../../scripts /etc/gm/scripts failed"
fi

ln -s $script_dir/conf/gm-master.toml /etc/gm/conf/gm.toml
if [ $? == 0 ]; then 
    echo "ln -s $script_dir/conf/gm-master.toml /etc/gm/conf/gm.toml succeed"
else
    echo "ln -s $script_dir/conf/gm-master.toml /etc/gm/conf/gm.toml failed"
fi

exit $?
